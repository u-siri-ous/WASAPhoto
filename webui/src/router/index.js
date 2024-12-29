import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import LogoutView from '../views/LogoutView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView, redirect:'/login'},
		{path: '/login', component: LoginView},
		{path: '/logout', component: LogoutView},
		{path: '/users/:userId', component: ProfileView},
	]
})

export default router
