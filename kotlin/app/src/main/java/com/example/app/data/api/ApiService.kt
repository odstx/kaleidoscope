package com.example.app.data.api

import com.example.app.data.model.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {
    @POST("api/v1/users/login")
    suspend fun login(@Body request: LoginRequest): LoginResponse

    @POST("api/v1/users/register")
    suspend fun register(@Body request: RegisterRequest): RegisterResponse

    @GET("api/v1/users/info")
    suspend fun getCurrentUser(): User

    @PUT("api/v1/users/info")
    suspend fun updateUser(@Body updates: Map<String, String?>): User

    @GET("api/v1/system/info")
    suspend fun getSystemInfo(): SystemInfo

    companion object {
        const val BASE_URL = "http://10.0.2.2:3000/"  // Android emulator localhost
    }
}

data class SystemInfo(
    val version: String,
    val status: String
) {
    constructor() : this("", "")
}
