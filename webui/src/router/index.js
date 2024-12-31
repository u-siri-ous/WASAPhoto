import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import LogoutView from '../views/LogoutView.vue'
import UserView from '../views/UserView.vue'
import ExploreUsers from '../views/ExploreUsers.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView, redirect:'/login'},
		{path: '/login', component: LoginView},
		{path: '/logout', component: LogoutView},
		{path: '/users/', component: ExploreUsers},
		{path: '/users/:userId', component: UserView},
	]
})

export default router
