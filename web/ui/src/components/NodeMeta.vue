<template>
  <b-container class="node-meta p-3">
    <b-row cols="4" align-content="start">
      <b-col align-self="start">
        <div class="info">
          Added by:
          <router-link
            :to="{ name: 'profile', params: { username: node.username } }"
            class="node"
            >
            {{ node.username }}
          </router-link><br />
            <span class="date">{{ node.CreatedAt | date }}</span>
        </div>
      </b-col>
    </b-row>
     <IotctlNodeActions
      v-if='actions'
      :node='node'
      :canModify='isCurrentUser()'
    />
  </b-container>
</template>

<script>
import { mapGetters } from 'vuex'
import IotctlNodeActions from '@/components/NodeActions'

export default {
  name: 'IotctlNodeMeta',
  components: {
    IotctlNodeActions
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
      return true
    }
  }
}

</script>

<style lang="scss">
.node-meta {
  text-align: left;
  width: 100%;
}
</style>
