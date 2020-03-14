<template>
  <div class="node-meta">
    <router-link
      :to="{ name: 'profile', params: { username: node.username } }"
    >

    </router-link>

    <div class="info">
      <router-link
        :to="{ name: 'profile', params: { username: node.username } }"
        class="node"
      >
        {{ node.username }}
      </router-link>
      <span class="date">{{ node.CreatedAt | date }}</span>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'IotctlNodeMeta',
  components: {

  },

  props: {
    node: {
      type: Object,
      required: true
    },
    actions: {
      type: Boolean,
      required: false,
      default: false
    }
  },

  computed: {
    ...mapGetters(['currentUser', 'isAuthenticated'])
  },

  methods: {
    isCurrentUser () {
      if (this.currentUser.username && this.node.username) {
        return this.currentUser.username === this.node.username
      }
      return false
    }
  }
}

</script>
