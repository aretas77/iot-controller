<template>
  <b-container>
    <b-row fluid class="m-3">

      <b-col>
        <p v-text="unregNode.mac" />
      </b-col>

    </b-row>
  </b-container>
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
