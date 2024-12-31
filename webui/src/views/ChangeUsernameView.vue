<template>
    <div class="change-username-container w-75 mx-auto my-5 p-4">
        <form @submit.prevent="submitForm" class="change-username-form">
            <h2 class="mb-4 text-center">Change Username</h2>
            
            <div class="mb-3">
                <label for="inputName" class="form-label">New Username</label>
                <input 
                    type="text" 
                    class="form-control" 
                    id="newUsername" 
                    aria-describedby="usernameHelp"
                    v-model="newUsername" 
                    :class="{ 'is-invalid': !isUsernameValid() }" 
                    placeholder="Enter your new username">
                <div class="invalid-feedback">
                    Username must be 3-25 characters long, start with a letter, and can contain letters, numbers, or underscores.
                </div>
            </div>
            
            <button 
                type="submit" 
                class="btn btn-primary w-100" 
                :disabled="!newUsername || !isUsernameValid()">Submit
            </button>

            <div 
                class="alert alert-success mt-3" 
                role="alert" 
                v-if="changeUsernameSuccess">
                Username changed successfully!
            </div>

            <div 
                class="alert alert-danger mt-3" 
                role="alert" 
                v-else-if="changeUsernameError">
                {{ error_msg }}
            </div>
        </form>
    </div>
</template>

<script>
const token = sessionStorage.getItem('authToken');
export default {
    data() {
        return {
            newUsername: '',
            changeUsernameSuccess: false,
            changeUsernameError: false,
            error_msg: '',
        };
    },
    methods: {
        async submitForm() {
            try {
                const config = {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`,
                    },
                };
                const response = await this.$axios.put(`/settings/username`, { username: this.newUsername }, config);
                this.changeUsernameSuccess = true;
                this.changeUsernameError = false;
            }
            catch (error) {
                console.error(error, "Error in changing name");
                const statusCode = error.response.status;
                switch (statusCode) {
                    case 401:
                        console.error('Access Unauthorized:', error.response.data);
                        this.error_msg = "You are not logged in"
                        break;
                    case 400:
                        console.error('Bad request:', error.response.data);
                        this.error_msg = "Name already in use"
                        break;
                    case 403:
                        console.error('Forbidden Action: ', error.response.data);
                        this.error_msg = "Username already in use"
                        break
                    case 404:
                        console.error('Not found: ', error.response.data);
                        this.error_msg = "You are not logged in"
                    default:
                        console.error(`Unhandled HTTP Error (${statusCode}):`, error.response.data);
                        this.error_msg = "You should login first"
                }
                this.changeUsernameSuccess = false;
                this.changeUsernameError = true;
            }

            this.newUsername = '';
        },
        isUsernameValid() {
			const usernameRegex = /^[a-zA-Z][\.]{0,1}([\w][\.]{0,1})*[\w]$/
			return usernameRegex.test(this.newUsername) && this.newUsername.length >= 5 && this.newUsername.length <= 25
		},
    },
};
</script>

<style scoped>
.change-username-container {
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #f9f9f9;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.change-username-form {
    text-align: center;
}

.change-username-form .form-label {
    font-weight: bold;
    text-align: left;
}

.btn-primary {
    font-size: 1rem;
    padding: 10px 20px;
    transition: background-color 0.3s ease, transform 0.2s ease;
}

.btn-primary:hover {
    background-color: #0056b3;
    transform: translateY(-2px);
}

.alert {
    text-align: center;
    font-size: 0.9rem;
}
</style>