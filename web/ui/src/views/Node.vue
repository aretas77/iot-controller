<template>
  <div class="node-page">
    <div class="banner">
      <div class="container">
        <h1>{{ node.name }}</h1>
      </div>
    </div>
    <div class="container page">
      <div class="row node-content">
      </div>
    </div>
    <hr />
    <div class='node-actions'>
      <NodeMeta :node='node' :actions='true'></NodeMeta>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import store from '@/store'
import marked from 'marked'
import NodeMeta from '@/components/NodeMeta'
import { FETCH_NODE } from '@/store/actions.type'

export default {
  name: 'iotctl-node',
  props: {
    slug: {
      type: Number,
      required: true
    }
  },
  components: {
    NodeMeta
  },
  // actions
  beforeRouteEnter (to, from, next) {
    console.log(to.params)
    Promise.all([
      store.dispatch(FETCH_NODE, to.params.slug)
    ]).then(() => {
      next()
    })
  },
  computed: {
    ...mapGetters(['node', 'currentUser', 'isAuthenticated'])
  },
  methods: {
    parseMarkdown (content) {
      return marked(content)
    }
  }
}

</script>
