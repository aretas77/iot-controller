<template>
  <div class="node-page mt-3">
    <div class="banner">
      <div class="container">
        <h1>{{ node.name }}</h1>
      </div>
    </div>
    <div>
      <b-tabs content-class="mt-3">
        <b-tab title="Statistics" active>
          <div class="container page">
            <div class="row node-content">
            </div>
          </div>
        </b-tab>

        <b-tab title="Models">

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
import { FETCH_NODE } from '@/store/actions.type'

export default {
  name: 'iotctl-node',
  props: {
    slug: {
      type: [Number, String],
      required: true
    }
  },
  components: {
    NodeMeta
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
    ...mapGetters(['node', 'currentUser', 'isAuthenticated'])
  },
  methods: {
    parseMarkdown (content) {
      return marked(content)
    }
  }
}
</script>
