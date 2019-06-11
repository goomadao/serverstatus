import Vue from 'vue'
import Router from 'vue-router'
// import Saved from '@/components/Saved.vue'
// import Temp from '@/components/Temp.vue'

Vue.use(Router)

function loadView (view) {
  return () => import(`@/components/${view}.vue`)
}

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Saved',
      component: loadView('Saved')
    },
    {
      path: '/temp',
      name: 'Temp',
      component: loadView('Temp')
    }
  ]
})
