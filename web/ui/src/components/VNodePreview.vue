<template>
  <div class="node-preview">
    <NodeMeta isPreview :node="node" />
    <router-link :to="nodeLink" class="preview-link">
      <h1 v-text="node.name" />
      <p v-text="node.mac" />
      <span>Read more...</span>
    </router-link>
  </div>
</template>

<script>
import NodeMeta from './NodeMeta'

export default {
  name: 'IotctlNodePreview',
  components: {
    NodeMeta
  },
  props: {
    node: { type: Object, required: true }
  },
  computed: {
    nodeLink () {
      return {
        name: 'node',
        params: {
          slug: this.node.ID
        }
      }
    }
  }
}
</script>
