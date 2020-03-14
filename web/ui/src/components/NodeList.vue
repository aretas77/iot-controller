<template>
  <div>
    <div v-if="isLoading" class="node-preview">Loading nodes...</div>
    <div v-else>
      <div v-if="nodes.length === 0" class="node-preview">
        No nodes are added.
      </div>
      <IotctlNodePreview
        v-for="(node, index) in nodes"
        :node="node"
        :key="node.mac + index"
      />
      <VPagination :pages="pages" :currentPage.sync="currentPage" />
    </div>
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
      default: 10
    }
  },
  data () {
    return {
      currentPage: 1
    }
  },
  computed: {
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
  methods: {
    fetchNodes () {
      this.$store.dispatch(FETCH_NODES, this.listConfig)
    },
    resetPagination () {
      this.currentPage = 1
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
    }
  },
  mounted () {
    this.fetchNodes()
  }
}
</script>
