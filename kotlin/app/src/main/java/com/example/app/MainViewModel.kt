package com.example.app

import androidx.lifecycle.ViewModel
import com.example.app.utils.TokenManager
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import javax.inject.Inject

@HiltViewModel
class MainViewModel @Inject constructor(
    private val tokenManager: TokenManager
) : ViewModel() {
    
    private val _isAuthenticated = MutableStateFlow(tokenManager.hasToken())
    val isAuthenticated: StateFlow<Boolean> = _isAuthenticated.asStateFlow()
    
    fun logout() {
        tokenManager.clearToken()
        _isAuthenticated.value = false
    }
}
