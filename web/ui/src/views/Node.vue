<template>
  <div class="node-page mt-3">
    <div class="banner">
      <div class="container">
        <h1>{{ node.name }}</h1>
        <p>{{ node.mac }}</p>
      </div>
    </div>
    <div>
      <b-tabs content-class="mt-3" @input="onChangedTab">

        <b-tab title="Statistics" active lazy>
          <b-container v-if="isLoadingStats">
            Loading...
          </b-container>

          <b-container fluid class="w-100 p-3" v-else>
            <b-row align-h="between" class="m-4">
              <b-col cols="5">
                <TemperatureChart :entries='statsEntries'></TemperatureChart>
              </b-col>
              <b-col cols="5">
                <SensorReadingsFreq />
              </b-col>
            </b-row>
          </b-container>

        </b-tab>

        <b-tab title="Models" lazy>

        </b-tab>

        <b-tab title="Information">
          <div class='node-actions'>
            <NodeMeta :node='node' :actions='true'></NodeMeta>
          </div>
        </b-tab>
      </b-tabs>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import store from '@/store'
import marked from 'marked'
import NodeMeta from '@/components/NodeMeta'
import TemperatureChart from '@/components/TemperatureChart'
import SensorReadingsFreq from '@/components/SensorReadingsFreq'
import { FETCH_NODE, FETCH_NODE_STATS } from '@/store/actions.type'

export default {
  name: 'iotctl-node',
  components: {
    NodeMeta,
    TemperatureChart,
    SensorReadingsFreq
  },
  props: {
    slug: {
      type: [Number, String],
      required: true
    }
  },
  data () {
    return {
      activeTab: 0
    }
  },
  // actions
  beforeRouteEnter (to, from, next) {
    Promise.all([
      store.dispatch(FETCH_NODE, to.params.slug)
    ]).then(() => {
      next()
    })
  },
  computed: {
    ...mapGetters(['node', 'currentUser', 'isAuthenticated', 'isLoadingStats',
      'statsEntries'])
  },
  methods: {
    parseMarkdown (content) {
      return marked(content)
    },
    fetchNodeStatistics () {
      this.$store.dispatch(FETCH_NODE_STATS, this.node.ID)
    },
    onChangedTab (tabIndex) {
      if (tabIndex === 0) {
        // refresh statistics
        // TODO: could probably just check with the server if there are new
        // updates
        this.$store.dispatch(FETCH_NODE_STATS, this.node.ID)
      }
    }
  },
  mounted () {
    this.fetchNodeStatistics()
  }
}
</script>
