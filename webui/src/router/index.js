import {createRouter, createWebHashHistory} from 'vue-router'
import HomeWasaText from '../views/HomeWasaText.vue'
import ProfileView from '../views/ProfileView.vue'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'
import LachiHome from '../views/LachiHome.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/home', component: HomeWasaText},
		{path: '/profiles/user', component: ProfileView},
		{path: '/home/:destUserID', component: ChatView},
		
		
	]
})

export default router
