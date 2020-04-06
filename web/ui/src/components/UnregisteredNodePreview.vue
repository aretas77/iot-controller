<template>
  <b-list-group-item class="m-2 text-left">
    <span class="unreg-mac">
      {{ unregNode.mac }}
    </span>
    <br />
    <span>
      Added at: {{ unregNode.CreatedAt | time }}
    </span>
    <br />
    <span>
      Added by:
      <router-link
        :to="{ name: 'profile', params: { username: unregNode.username } }"
        >
        {{ unregNode.username }}
      </router-link>
    </span>
  </b-list-group-item>
</template>

<script>

export default {
  name: 'UnregisteredNodePreview',
  props: {
    unregNode: {
      type: Object,
      required: true
    }
  }
}
</script>

<style lang="scss">
.unreg-mac {
  color: #5cb85c;
  text: bold;
}
</style>
