package com.example.app.ui.screens

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.app.data.model.UpdateUserRequest
import com.example.app.data.repository.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ProfileUiState(
    val isLoading: Boolean = false,
    val username: String = "",
    val email: String = "",
    val isUpdating: Boolean = false,
    val updateError: String? = null,
    val isUpdateSuccess: Boolean = false,
    val errorMessage: String? = null
)

@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {
    
    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()
    
    fun loadProfile() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, errorMessage = null)
            
            val result = authRepository.getUserInfo()
            
            result.fold(
                onSuccess = { userInfo ->
                    _uiState.value = ProfileUiState(
                        isLoading = false,
                        username = userInfo.username,
                        email = userInfo.email
                    )
                },
                onFailure = { error ->
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        errorMessage = error.message ?: "Failed to load profile"
                    )
                }
            )
        }
    }
    
    fun updateUsername(username: String) {
        _uiState.value = _uiState.value.copy(username = username, updateError = null)
    }
    
    fun updateEmail(email: String) {
        _uiState.value = _uiState.value.copy(email = email, updateError = null)
    }
    
    fun updateProfile() {
        val currentState = _uiState.value
        
        if (currentState.username.isBlank()) {
            _uiState.value = currentState.copy(updateError = "Username is required")
            return
        }
        
        if (currentState.email.isBlank()) {
            _uiState.value = currentState.copy(updateError = "Email is required")
            return
        }
        
        if (!android.util.Patterns.EMAIL_ADDRESS.matcher(currentState.email).matches()) {
            _uiState.value = currentState.copy(updateError = "Invalid email address")
            return
        }
        
        viewModelScope.launch {
            _uiState.value = currentState.copy(isUpdating = true, updateError = null)
            
            val result = authRepository.updateUserInfo(
                UpdateUserRequest(
                    username = currentState.username,
                    email = currentState.email
                )
            )
            
            result.fold(
                onSuccess = {
                    _uiState.value = currentState.copy(
                        isUpdating = false,
                        isUpdateSuccess = true
                    )
                },
                onFailure = { error ->
                    _uiState.value = currentState.copy(
                        isUpdating = false,
                        updateError = error.message ?: "Failed to update profile"
                    )
                }
            )
        }
    }
}
