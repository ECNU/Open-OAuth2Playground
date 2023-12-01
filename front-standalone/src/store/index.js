import { createStore } from "vuex"
import createPersistedState from 'vuex-persistedstate';


const store = createStore({
  state() {
    return {
      backendUrl: 'http://localhost:80',
      authServerUrl: '',
      grantTypes: "1"
    }
  },
  mutations: {
    setBackEndUrl(state, v) {
      state.backendUrl = v
    },
    setGrantTypes(state, v) {
      state.grantTypes = v;
    },
  },
  plugins: [createPersistedState()],
})

export default store
