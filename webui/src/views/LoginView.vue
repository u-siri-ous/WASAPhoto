<template>
    <div class="login-container">
        <LoadingSpinner v-if="loading"></LoadingSpinner>
        <div class="login-form">
            <h2>Login</h2>
            <form @submit.prevent="login">
                <label class="login-label" for="username">Username:</label>
                <input type="text" class="form-control" id="username" aria-describedby="usernameHelp" v-model="username"
                    :class="{ 'is-invalid': !isUsernameValid() }">
                <button type="submit" class="btn btn-sm btn-outline-primary login-btn"
                    :disabled="!username || !isUsernameValid() || loading"
                    style="align-self: center; float: center; font-size: 20px;">Login <svg class="feather">
                        <use href="/feather-sprite-v4.29.0.svg#key" />
                    </svg></button>
            </form>
            <div v-if="identity !== null">
                <p>User logged successfully with id: {{ identity.userId }}</p>
            </div>
        </div>
    </div>
</template>

<script>
import LoadingSpinner from '@/components/LoadingSpinner.vue';

const token = sessionStorage.getItem('authToken');
export default {
    data() {
        return {
            username: "",
            identity: null,
            loading: false,
        };
    },
    components: {
        LoadingSpinner,
    },
    methods: {
        async login() {
            this.loading = true;
            this.errormsg = null;
            console.log(this.username)
            try {
                let response = await this.$axios.post('/login', { username: this.username }, {
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json',
                    },
                });
                console.log(response)
                this.identity = response.data
                this.SaveToSessionStorage()
            } catch (error) {
                console.error("Error while logging in!");
            }
            this.loading = false;
            this.navigateToMyPage()
        },
        navigateToMyPage() {
            localStorage.setItem('loggedIn', '1')
            location.reload();
        },
        SaveToSessionStorage() {
            const bearerToken = `${this.identity.userId}`;
            sessionStorage.setItem('authToken', bearerToken);
            sessionStorage.setItem('username', this.identity.username);
            sessionStorage.setItem('userId', this.identity.userId);
        },
        isUsernameValid() {

            const usernameRegex = /^[a-zA-Z][\.]{0,1}([\w][\.]{0,1})*[\w]$/
            return usernameRegex.test(this.username) && this.username.length >= 5 && this.username.length <= 25
        },
    },
};
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    text-align: center;
    align-items: center;
    height: 100vh;
}

.login-form {
    padding: 20px;
    border: none;
    align-items: center;
    border-radius: 8px;
}

.login-label {
    padding: 3px;
    display: block;
    margin-bottom: 8px;
}

.login-btn {
    align-self: center;
}
</style>