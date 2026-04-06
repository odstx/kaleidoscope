package com.example.app.data.repository

import com.example.app.data.model.*
import com.example.app.data.api.ApiService
import com.example.app.utils.TokenManager
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val apiService: ApiService,
    private val tokenManager: TokenManager
) {
    private val _currentUser = MutableStateFlow<User?>(null)
    val currentUser: StateFlow<User?> = _currentUser.asStateFlow()

    private val _isAuthenticated = MutableStateFlow(false)
    val isAuthenticated: StateFlow<Boolean> = _isAuthenticated.asStateFlow()

    private val _isLoading = MutableStateFlow(false)
    val isLoading: StateFlow<Boolean> = _isLoading.asStateFlow()

    private val _error = MutableStateFlow<String?>(null)
    val error: StateFlow<String?> = _error.asStateFlow()

    suspend fun login(email: String, password: String): Result<LoginResponse> {
        _isLoading.value = true
        _error.value = null
        
        return try {
            val response = apiService.login(LoginRequest(email, password))
            tokenManager.saveToken(response.token)
            _currentUser.value = response.user
            _isAuthenticated.value = true
            _isLoading.value = false
            Result.success(response)
        } catch (e: Exception) {
            _isLoading.value = false
            _error.value = e.message ?: "Login failed"
            Result.failure(e)
        }
    }

    suspend fun register(request: RegisterRequest): Result<RegisterResponse> {
        _isLoading.value = true
        _error.value = null
        
        return try {
            val response = apiService.register(request)
            tokenManager.saveToken(response.token)
            _currentUser.value = response.user
            _isAuthenticated.value = true
            _isLoading.value = false
            Result.success(response)
        } catch (e: Exception) {
            _isLoading.value = false
            _error.value = e.message ?: "Registration failed"
            Result.failure(e)
        }
    }

    suspend fun loadCurrentUser(): Result<User> {
        _isLoading.value = true
        _error.value = null
        
        return try {
            val user = apiService.getCurrentUser()
            _currentUser.value = user
            _isAuthenticated.value = true
            _isLoading.value = false
            Result.success(user)
        } catch (e: Exception) {
            _isLoading.value = false
            _error.value = e.message ?: "Failed to load user"
            Result.failure(e)
        }
    }

    suspend fun updateUser(name: String?, bio: String?): Result<User> {
        _isLoading.value = true
        _error.value = null
        
        return try {
            val user = apiService.updateUser(mapOf(
                "name" to name,
                "bio" to bio
            ))
            _currentUser.value = user
            _isLoading.value = false
            Result.success(user)
        } catch (e: Exception) {
            _isLoading.value = false
            _error.value = e.message ?: "Failed to update user"
            Result.failure(e)
        }
    }

    fun logout() {
        tokenManager.clearToken()
        _currentUser.value = null
        _isAuthenticated.value = false
    }

    fun clearError() {
        _error.value = null
    }

    fun checkAuthState(): Boolean {
        val hasToken = tokenManager.hasToken()
        _isAuthenticated.value = hasToken
        return hasToken
    }
    
    suspend fun getUserInfo(): Result<UserInfo> {
        return try {
            val user = apiService.getCurrentUser()
            Result.success(UserInfo(user.id, user.email, user.username))
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    suspend fun getSystemInfo(): Result<SystemInfo> {
        return try {
            val info = apiService.getSystemInfo()
            Result.success(info)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
    
    suspend fun updateUserInfo(request: UpdateUserRequest): Result<UserInfo> {
        return try {
            val user = apiService.updateUser(mapOf(
                "username" to request.username,
                "email" to request.email
            ))
            Result.success(UserInfo(user.id, user.email, user.username))
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
