<template>
  <div class="node-editor">
    <b-container>
      <b-row>
        <b-col>
          <IotctlListErrors :errors="errors" />
          <form @submit.prevent="onPublish(node.slug)">
            <fieldset :disabled="inProgress">
              <fieldset class="form-group">
                <input
                  type="text"
                  class="form-control form-control-lg"
                  v-model="node.name"
                  placeholder="Node Name"
                />
              </fieldset>
            </fieldset>
            <button
              :disabled="inProgress"
              class="btn btn-lg pull-xs-right btn-primary"
              type="submit"
            >
              Register Node
            </button>
          </form>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import store from '@/store'
import IotctlListErrors from '@/components/ListErrors'
import {
  NODE_ADD,
  NODE_EDIT,
  FETCH_NODE,
  NODE_RESET_STATE
} from '@/store/actions.type'

export default {
  name: 'NodeCreate',
  components: {
    IotctlListErrors
  },
  props: {
    previousNode: {
      type: Object,
      required: false
    }
  },

  async beforeRouterUpdate (to, from, next) {
    // Reset state if user goes from /editor/:id to /editor
    // The component is not recreated so we use hook to reset the state.
    await store.dispatch(NODE_RESET_STATE)
    return next()
  },

  async beforeRouteEnter (to, from, next) {
    await store.dispatch(NODE_RESET_STATE)
    if (to.params.slug !== undefined) {
      await store.dispatch(
        FETCH_NODE,
        to.params.slug,
        to.params.previousNode
      )
    }
    return next()
  },

  data () {
    return {
      inProgress: false,
      errors: {}
    }
  },
  computed: {
    ...mapGetters(['node'])
  },

  methods: {
    onPublish (slug) {
      const action = slug ? NODE_EDIT : NODE_ADD
      this.inProgress = true
      this.$store
        .dispatch(action)
        .then(({ data }) => {
          this.inProgress = false
          this.$router.push({
            name: 'node',
            params: { slug: data.node.slug }
          })
        }).catch(({ response }) => {
          this.inProgress = false
          this.errors = response.data.errors
        })
    }
  }
}
</script>
