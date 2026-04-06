package com.example.app.ui.screens

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.example.app.data.model.SystemInfo
import com.example.app.data.model.UserInfo
import com.example.app.data.repository.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class DashboardUiState(
    val isLoading: Boolean = false,
    val userInfo: UserInfo? = null,
    val systemInfo: SystemInfo? = null,
    val errorMessage: String? = null
)

@HiltViewModel
class DashboardViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {
    
    private val _uiState = MutableStateFlow(DashboardUiState())
    val uiState: StateFlow<DashboardUiState> = _uiState.asStateFlow()
    
    fun loadDashboardData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, errorMessage = null)
            
            val userInfoResult = authRepository.getUserInfo()
            val systemInfoResult = authRepository.getSystemInfo()
            
            val userInfo = userInfoResult.getOrNull()
            val systemInfo = systemInfoResult.getOrNull()
            
            val errorMessage = when {
                userInfoResult.isFailure -> userInfoResult.exceptionOrNull()?.message ?: "Failed to load user info"
                systemInfoResult.isFailure -> systemInfoResult.exceptionOrNull()?.message ?: "Failed to load system info"
                else -> null
            }
            
            _uiState.value = DashboardUiState(
                isLoading = false,
                userInfo = userInfo,
                systemInfo = systemInfo,
                errorMessage = errorMessage
            )
        }
    }
}
