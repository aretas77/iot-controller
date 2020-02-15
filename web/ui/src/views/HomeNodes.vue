<template>
  <div class="home-nodes">
    <NodeList type="all" />
  </div>
</template>

<script>
import NodeList from '@/components/NodeList'

export default {
  name: 'home-nodes',
  components: {
    NodeList
  }
}
</script>
