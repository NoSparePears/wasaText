<!-- 
	In questo file si trova la barra di navigazione (posta a sx)
-->

<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return {
			errorMsg: null,

			// Utilizato per mostarer determinati contenuti della pagina solo se un utente ha effettuato il login
			isLoggedIn: sessionStorage.token ? true : false,

			// UserId dell'utente loggato
			userID: sessionStorage.userID,

			// Username dell'utente loggato
			username: sessionStorage.username,
			// Profile picture dell'utente logagto
			photo: sessionStorage.photo,
			// Utilizzato per controllare se l'username inserito dall'utente è valido
			usernameValidation: new RegExp('^\\w{3,16}$'),
		}
	},

	methods: {
		// Funzione utilizzata per il login dell'utente
		handleLoginSuccess() {
      		this.isLoggedIn = true;
      		this.userID = sessionStorage.userID;
      		this.username = sessionStorage.username;
      		this.photo = sessionStorage.photo;
    	},
		// Funzione utilizzata per il logout dell'utente
		logOut() {
    	  sessionStorage.clear();
    	  this.isLoggedIn = false;
    	  this.$router.push("/");
    	}
	}
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">wasaText</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item" v-if="isLoggedIn">
							<RouterLink to="/home" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
								Chat
							</RouterLink>
						</li>
						<li class="nav-item" v-if="isLoggedIn">
							<RouterLink to="/profiles/user" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
								Profilo
							</RouterLink>
						</li>
						<!-- Esegue il logout (solo se l'utente è loggato) ritornando alla pagina di login -->
						<li class="nav-item" v-if="isLoggedIn">
							<a class="nav-link" @click="logOut">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
								Log Out
							</a>	
						</li>
					</ul>

				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @login-success="handleLoginSuccess"/>
			</main>
		</div>
	</div>
</template>

<!-- Stili per l'immagine del profilo e dell'username dell'utente loggato -->
<style>
.profile-picture {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.username {
  font-size: 14px;
  font-weight: bold;
}
</style>
