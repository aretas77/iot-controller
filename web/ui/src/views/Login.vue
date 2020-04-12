<template>
  <div class="auth-page">
    <b-container>
      <b-row>
        <b-col class="col-md-6 offset-md-3 col-xs-12 mt-5">
          <h1 class="text-xs-center">Sign in</h1>

            <b-list-group v-if="authErrors" class="error-messages m-3">
              <b-list-group-item
                v-for="(v, k) in authErrors"
                :key="k"
                variant="danger"
                >
                {{ v }}
              </b-list-group-item>
            </b-list-group>

          <form @submit.prevent="onSubmit(email, password)">
            <fieldset class="form-group">
              <input
                class="form-control form-control-lg"
                type="text"
                v-model="email"
                placeholder="Email"
              />
            </fieldset>
            <fieldset class="form-group">
              <input
                class="form-control form-control-lg"
                type="password"
                v-model="password"
                placeholder="Password"
              />
            </fieldset>

            <button class="btn btn-lg btn-primary pull-xs-right">
              Sign in
            </button>
          </form>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import { LOGIN } from '@/store/actions.type'

export default {
  name: 'Login',
  data () {
    return {
      email: null,
      password: null
    }
  },
  methods: {
    onSubmit (email, password) {
      this.$store
        .dispatch(LOGIN, { email, password })
        .then(() => this.$router.push({ name: 'Home' }))
    }
  },
  computed: {
    ...mapGetters(['authErrors'])
  }
}
</script>
