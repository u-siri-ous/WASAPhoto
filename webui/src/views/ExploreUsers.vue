<template>
    <div class="user-search-container">
        <h2 class="title">Explore Users</h2>
        <form @submit.prevent="searchUsers" class="search-form">
            <div class="input-container">
                <label for="searchQuery" class="search-label">Username:</label>
                <input 
                    type="text" 
                    id="searchQuery" 
                    v-model="searchQuery" 
                    class="search-input" 
                    placeholder="Enter username"
                    :class="{ 'is-invalid': !isUsernameValid() }"
                />
                <button 
                    type="submit" 
                    class="search-button" 
                    :disabled="!searchQuery || !isUsernameValid()"
                >
                    Search
                </button>
            </div>
        </form>

        <p v-if="searched" class="result-text">{{ Text }}</p>

        <div v-if="Users.length" class="user-grid">
            <div v-for="user in Users" :key="user.userId" class="user-card">
                <p class="user-name">{{ user.username }}</p>
                <button 
                    type="button" 
                    class="profile-button" 
                    @click="$router.push(`/users/${user.userId}`)"
                >
                    View Profile
                </button>
            </div>
        </div>
        <p v-else-if="searched" class="no-results">No users found.</p>
    </div>
</template>

<script>
const token = sessionStorage.getItem('authToken');

export default {
    data() {
        return {
            searchQuery: '',
            searched: false,
            Text: '',
            Users: [],
        };
    },
    methods: {
        async searchUsers() {
            try {
                const response = await this.$axios.get(`/users/`, {
                    params: { username: this.searchQuery },
                    headers: {
                        'Authorization': `Bearer ${token}`,
                        'Accept': 'application/json',
                    },
                });
                this.searched = true;
                this.Users = response.data.users;
                this.Text = this.Users.length === 0 ? "No users found matching your search." : "";
            } catch (error) {
                console.error(error, "Error during user search");
                this.searched = true;
                if (error.response) {
                    const statusCode = error.response.status;
                    switch (statusCode) {
                        case 401:
                            console.error('Unauthorized Access:', error.response.data);
                            this.Text = "You need to log in to perform a search.";
                            break;
                        case 403:
                            console.error('Forbidden Access:', error.response.data);
                            this.Text = "You don't have the necessary permissions to perform this search.";
                            break;
                        case 404:
                            console.error('Resource Not Found:', error.response.data);
                            this.Text = "No users found with the provided username.";
                            break;
                        default:
                            console.error(`Unexpected Error (${statusCode}):`, error.response.data);
                            this.Text = "An unexpected error occurred. Please try again later.";
                    }
                } else {
                    console.error('Network or Server Error:', error);
                    this.Text = "Unable to connect to the server. Please check your internet connection and try again.";
                }
            }
        },
        isUsernameValid() {
            const usernameRegex = /^[a-zA-Z][\.]{0,1}([\w][\.]{0,1})*$/;
            return usernameRegex.test(this.searchQuery) && this.searchQuery.length >= 3 && this.searchQuery.length <= 25;
        },
    },
};
</script>

<style scoped>
.user-search-container {
    margin: 0 auto;
    padding: 20px;
    font-family: 'Roboto', sans-serif;
    text-align: center;
}

.title {
    font-size: 2rem;
    margin-bottom: 1.5rem;
    color: #333;
}

.search-form {
    margin-bottom: 2rem;
}

.input-container {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
}

.search-label {
    font-size: 1.2rem;
    font-weight: bold;
    color: #555;
}

.search-input {
    flex: 1;
    max-width: 400px;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 4px;
    outline: none;
    font-size: 1rem;
    transition: border-color 0.3s ease;
}

.search-input.is-invalid {
    border-color: red;
}

.search-input:focus {
    border-color: #007bff;
}

.search-button {
    padding: 10px 20px;
    font-size: 1rem;
    border: none;
    border-radius: 4px;
    background-color: #007bff;
    color: #fff;
    cursor: pointer;
    transition: background-color 0.3s ease;
}

.search-button:disabled {
    background-color: #ccc;
    cursor: not-allowed;
}

.result-text {
    margin-top: 1rem;
    font-size: 1.2rem;
    color: #333;
}

.user-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 20px;
    margin-top: 2rem;
}

.user-card {
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 8px;
    background-color: #f9f9f9;
    text-align: center;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.user-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 12px rgba(0, 0, 0, 0.2);
}

.user-name {
    font-size: 1.2rem;
    font-weight: bold;
    margin-bottom: 10px;
    color: #333;
}

.profile-button {
    padding: 8px 16px;
    font-size: 1rem;
    border: none;
    border-radius: 4px;
    background-color: #28a745;
    color: #fff;
    cursor: pointer;
    transition: background-color 0.3s ease;
}

.profile-button:hover {
    background-color: #218838;
}

.no-results {
    font-size: 1.2rem;
    color: #999;
    margin-top: 1rem;
}
</style>