<template>
  <div class="node-page mt-3">
    <div class="banner">
      <div class="banner-info">
        <div class="py-4">
          <h1>{{ node.name }}</h1>
          <p class="mb-0">{{ node.mac }}</p>
        </div>
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
          <b-container v-if="isLoadingEvents" fluid class="w-100 p-3">
            <NodeEvents :events='events' :busy='isLoadingEvents' :eventsLen='eventsCount'>
            </NodeEvents>
          </b-container>
          <b-container fluid class="w-100 p-3" v-else>
            <NodeEvents :events='events' :busy='isLoadingEvents' :eventsLen='eventsCount'>
            </NodeEvents>
          </b-container>
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
import NodeEvents from '@/components/Events'
import {
  FETCH_NODE,
  FETCH_NODE_STATS,
  FETCH_NODE_EVENTS,
  CHECK_AUTH
} from '@/store/actions.type'

export default {
  name: 'iotctl-node',
  components: {
    NodeMeta,
    TemperatureChart,
    SensorReadingsFreq,
    NodeEvents
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
      'statsEntries', 'events', 'isLoadingEvents', 'eventsCount'])
  },
  methods: {
    parseMarkdown (content) {
      return marked(content)
    },
    fetchNodeStatistics () {
      this.$store.dispatch(FETCH_NODE_STATS, this.node.ID)
    },
    fetchNodeEvents () {
      this.$store.dispatch(FETCH_NODE_EVENTS, this.node.ID)
    },
    onChangedTab (tabIndex) {
      this.$store.dispatch(CHECK_AUTH)

      if (tabIndex === 0) {
        // refresh statistics
        // TODO: could probably just check with the server if there are new
        // updates
        this.fetchNodeStatistics()
      } else if (tabIndex === 1) {
        this.fetchNodeEvents()
      }
    }
  },
  mounted () {
    this.fetchNodeStatistics()
    this.fetchNodeEvents()
  }
}
</script>

<style lang="scss">
.banner {
  background: #5cb85c;
  box-shadow: inset 0 8px 8px -8px rgba(0,0,0,.3),inset 0 -8px 8px -8px rgba(0,0,0,.3);
  color: #fff;
  margin-bottom: 1rem;
  margin-top: 1rem;
}
.banner-info {
  text-shadow: 0 1px 3px rgba(0,0,0,.3);
  font-weight: 700 !important;

}
</style>
