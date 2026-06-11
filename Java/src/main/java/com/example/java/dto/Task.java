package com.example.java.dto;

import jakarta.validation.constraints.NotBlank;

public record Task(@NotBlank(message = "Key cannot be blank") String key,
                   @NotBlank(message = "Value cannot be blank") String value) {
}
