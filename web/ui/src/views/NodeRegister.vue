<template>
  <div class="node-editor">
    <b-container>
      <b-row>
        <b-col>
          <IotctlListErrors :errors="errors" />

          <b-form @submit.prevent="onPublish(unregisteredNode.slug)">

            <fieldset :disabled="inProgress">
              <fieldset class="form-group">
                <input
                  type="text"
                  class="form-control form-control-lg"
                  v-model="unregisteredNode.mac"
                  placeholder="Node MAC"
                />
              </fieldset>

              <b-form-group>
                <b-form-select
                  v-model="network"
                  :options="networksOptions"
                  >
                  <template v-slot:first v-if="currentNetwork.name">
                    <b-form-select-option
                      :value="null"
                      disabled
                      >
                      -- Current network: {{ currentNetwork.name }} --
                    </b-form-select-option>
                  </template>
                </b-form-select>
              </b-form-group>
            </fieldset>

            <!-- Form submit button -->
            <button
              :disabled="inProgress"
              class="btn btn-lg pull-xs-right btn-primary"
              type="submit"
            >
              Register Node
            </button>
          </b-form>
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
  NODE_RESET_STATE,
  FETCH_NETWORKS
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
  mounted () {
    this.fetchNetworks()
  },
  watch: {
    // If there was a change in `networks` we build a new option list.
    networks (newValue, oldValue) {
      for (const key of newValue) {
        this.networksOptions.push({ text: key.name, value: key })
      }
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
    }
    return next()
  },

  data () {
    return {
      inProgress: false,
      errors: {},
      unregisteredNode: {},
      network: null,
      networksOptions: []
    }
  },
  computed: {
    ...mapGetters(['currentUser', 'currentNetwork', 'networks'])
  },

  methods: {
    // Publish the filled up form of Node creation
    onPublish (slug) {
      const action = NODE_ADD
      this.inProgress = true
      this.$store
        .dispatch(action)
        .then(({ data }) => {
          this.inProgress = false
          this.$router.push({
            name: 'node',
            params: { slug: data.unregisteredNode.slug }
          })
        }).catch(({ response }) => {
          this.inProgress = false
          // this.errors = response.data.errors
        })
    },
    fetchNetworks () {
      this.$store.dispatch(FETCH_NETWORKS, this.currentUser.ID)
    }
  }
}
</script>
