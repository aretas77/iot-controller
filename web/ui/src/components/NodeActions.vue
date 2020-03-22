<template>
  <!-- Used when User is also owner of the Node or Admin -->
  <span v-if="canModify">
    <router-link class="btn btn-sm btn-outline-secondary mr-3" :to="editNodeLink">
      <i class="ion-edit"></i> <span>&nbsp;Edit Node</span>
    </router-link>
    <span>&nbsp;&nbsp;</span>
    <button class="btn btn-outline-danger btn-sm" @click="deleteNode">
      <i class="ion-trash-a"></i> <span>&nbsp;Delete Node</span>
    </button>
  </span>
</template>

<script>
import { mapGetters } from 'vuex'
import {
  NODE_REMOVE
} from '@/store/actions.type'

export default {
  name: 'IotctlNodeActions',
  props: {
    node: {
      type: Object,
      required: true
    },
    canModify: {
      type: Boolean,
      required: true
    }
  },

  computed: {
    ...mapGetters(['profile', 'isAuthenticated']),
    editNodeLink () {
      // return { name: 'node-edit', params: { slug: this.node.slug } }
      return { name: 'Home' }
    }
  },

  methods: {
    async deleteNode () {
      try {
        await this.$store.dispatch(NODE_REMOVE, this.node.ID)
        this.$router.push('/')
      } catch (err) {
        console.error(err)
      }
    }
  }
}

</script>
