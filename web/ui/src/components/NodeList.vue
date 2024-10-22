<template>
  <div>
    <b-container v-if="isLoading" class="mh-100 mt-3">
      <b-progress-bar
        label="Loading nodes..."
        :value="value"
        :max=100
        animated>
      </b-progress-bar>
    </b-container>

    <b-container fluid class="node-list-container" v-else>
      <!-- Top of container for buttons and various editing -->
      <b-row class="m-3">
        <b-col col class="d-flex align-items-start">
          <router-link
            active-class="active"
            :to="{ name: 'node-register' }"
            >
            <b-button>
              Register node
            </b-button>
          </router-link>
        </b-col>

        <b-col>
          <VPagination :pages="pages" :currentPage.sync="currentPage" />
        </b-col>
      </b-row>

      <!-- Start of NodeList -->
      <b-row>
        <b-col>
          <div v-if="nodes.length === 0" class="node-preview">
            No nodes are added.
          </div>
          <IotctlNodePreview
            v-for="(node, index) in nodes"
            :node="node"
            :key="node.mac + index"
            />
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import IotctlNodePreview from './VNodePreview.vue'
import VPagination from './VPagination.vue'
import { FETCH_NODES } from '../store/actions.type'

export default {
  name: 'NodeList',
  components: {
    IotctlNodePreview,
    VPagination
  },
  props: {
    type: {
      type: String,
      required: false,
      default: 'all'
    },
    id: {
      type: Number,
      required: false
    },
    itemsPerPage: {
      type: Number,
      required: false,
      default: 5
    }
  },
  data () {
    return {
      value: 0,
      currentPage: 1,
      interval: null
    }
  },
  computed: {
    nodeCreate () {
      return {
        name: 'createNode'
      }
    },
    listConfig () {
      const { type } = this
      const filters = {
        offset: (this.currentPage - 1) * this.itemsPerPage,
        limit: this.itemsPerPage
      }
      if (this.id) {
        filters.id = this.id
      }

      return {
        type,
        filters
      }
    },
    pages () {
      if (this.isLoading || this.nodesCount <= this.itemsPerPage) {
        return []
      }
      return [
        ...Array(Math.ceil(this.nodesCount / this.itemsPerPage)).keys()
      ].map(e => e + 1)
    },
    ...mapGetters(['nodesCount', 'isLoading', 'nodes'])
  },
  created () {
    // this.pollNodeList()
  },
  beforeDestroy () {
    clearInterval(this.interval)
  },
  methods: {
    fetchNodes () {
      this.value = 50
      this.$store.dispatch(FETCH_NODES, this.listConfig)
    },
    resetPagination () {
      this.currentPage = 1
    },
    pollNodeList () {
      this.interval = setInterval(() => {
        this.fetchNodes()
      }, 2000)
    }
  },
  watch: {
    currentPage (newValue) {
      this.listConfig.filters.offset = (newValue - 1) * this.itemsPerPage
      this.fetchNodes()
    },
    type () {
      this.resetPagination()
      this.fetchNodes()
    },
    $route: 'fetchNodes'
  },
  mounted () {
    this.value = 0
    this.fetchNodes()
    this.value = 100
  }
}
</script>

<style lang="scss">
.node-list-container {

}
</style>
