import { createStore } from "vuex"


const store = createStore({
  state() {
    return {
      backendUrl: 'http://localhost:80',
      authServerUrl: ''
    }
  },
  mutations: {
    setBackEndUrl(state, v) {
      state.backendUrl = v
    }
  }
})


export default store