package com.example.app.data.model

import com.google.gson.annotations.SerializedName

data class User(
    val id: Long,
    val email: String,
    val username: String?,
    val name: String?,
    val bio: String?,
    val avatar: String?,
    @SerializedName("created_at")
    val createdAt: String?,
    @SerializedName("updated_at")
    val updatedAt: String?
)

data class LoginRequest(
    val email: String,
    val password: String
)

data class LoginResponse(
    val token: String,
    val user: User
)

data class RegisterRequest(
    val email: String,
    val password: String,
    val username: String?,
    val name: String?
)

data class RegisterResponse(
    val token: String,
    val user: User
)

data class ApiError(
    val error: String?,
    val message: String?,
    @SerializedName("error_description")
    val errorDescription: String?
)

data class UserInfo(
    val id: Long,
    val email: String,
    val username: String?
)

data class UpdateUserRequest(
    val username: String?,
    val email: String?
)
