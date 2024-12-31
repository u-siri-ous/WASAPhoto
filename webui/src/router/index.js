import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import LogoutView from '../views/LogoutView.vue'
import UserView from '../views/UserView.vue'
import ExploreUsers from '../views/ExploreUsers.vue'
import ChangeUsernameView from '../views/ChangeUsernameView.vue'
import CreatePostView from '../views/CreatePostView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView, redirect:'/login'},
		{path: '/login', component: LoginView},
		{path: '/logout', component: LogoutView},
		{path: '/users/', component: ExploreUsers},
		{path: '/settings', component: ChangeUsernameView},
		{path: '/posts', component: CreatePostView},
		{path: '/users/:userId', component: UserView},
	]
})

export default router