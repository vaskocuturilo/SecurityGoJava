package com.example.java.dto;

import com.example.java.entity.UserRole;
import lombok.Builder;

import java.util.Set;

@Builder
public record UserDto(String id, String username, String email, boolean active, Set<UserRole> roles) {
}

