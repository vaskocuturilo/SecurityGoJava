package com.example.java.dto;

import lombok.Builder;

@Builder
public record UserDto(String id, String username, String email, boolean active){}

