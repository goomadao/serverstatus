import Vue from 'vue'
import Router from 'vue-router'
import Saved from '@/components/Saved.vue'
import Temp from '@/components/Temp.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Saved',
      component: Saved
    },
    {
      path: '/temp',
      name: 'Temp',
      component: Temp
    }
  ]
})
