<template>
  <b-container fluid class="node-preview-container">
    <b-row class="m-3 node-preview-row ">

      <b-col cols=8 class="d-flex align-items-start">
        <router-link :to="nodeLink" class="preview-link">
          <h3 v-text="node.name" />
          <p v-text="node.mac" />
        </router-link>
      </b-col>

      <b-col cols=4>
        <NodeMeta isPreview :node="node" />
      </b-col>

    </b-row>
    <hr/>
  </b-container>
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

<style lang="scss">
.node-preview-container {
  margin: 7px;
}

.node-preview-row {
}

.preview-link {
}
</style>
