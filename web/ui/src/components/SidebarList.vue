<template>
  <b-container>
    <b-button v-b-toggle.sidebar-right @click="fetchUnregisteredNodes">&laquo;</b-button>
    <b-sidebar
      id="sidebar-right"
      title="Pending nodes"
      header-class="header-nodes2"
      right shadow
      >
      <hr />
      <b-list-group>
        <UnregisteredNodePreview
          v-for="(unode, index) in unregisteredNodes"
          :key="unode.mac + index"
          :unregNode="unode"
          />
      </b-list-group>
    </b-sidebar>
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import UnregisteredNodePreview from './UnregisteredNodePreview.vue'
import { FETCH_UNREG_NODES } from '@/store/actions.type'

export default {
  name: 'SidebarList',
  components: {
    UnregisteredNodePreview
  },
  data () {
    return { }
  },
  props: {
    networkName: {
      type: String,
      required: false,
      default: 'global'
    }
  },
  computed: {
    ...mapGetters(['unregisteredNodes', 'unregisteredNodesCount'])
  },
  methods: {
    fetchUnregisteredNodes () {
      this.$store.dispatch(FETCH_UNREG_NODES, this.networkName)
    }
  },
  mounted () {
    this.fetchUnregisteredNodes()
  }
}
</script>

<style lang="scss">
.header-nodes1 {
  background: #5cb85c;
  box-shadow: inset 0 8px 8px -8px rgba(0,0,0,.3),inset 0 -8px 8px -8px rgba(0,0,0,.3);
  color: #fff;
  margin-top: 1rem;
}
.header-nodes2 {
  color: #5cb85c;
}
</style>
