package com.example.java.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.RequiredArgsConstructor;

@Builder
@Data
@AllArgsConstructor
@RequiredArgsConstructor
public class UserDto {
    private String id;
    private String username;
    private String email;
    private String accessToken;
    private boolean active;
}
