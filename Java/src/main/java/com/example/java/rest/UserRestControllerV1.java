package com.example.java.rest;


import com.example.java.dto.*;
import com.example.java.exception.UserException;
import com.example.java.service.IUserService;
import com.example.java.service.RefreshTokenService;
import jakarta.validation.Valid;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.net.URI;

@RequiredArgsConstructor
@RestController
@RequestMapping("/api/v1/users")
public class UserRestControllerV1 {

    private final IUserService userService;
    private final RefreshTokenService refreshTokenService;

    @PostMapping("/login")
    public ResponseEntity<@NonNull AuthResponse> login(@RequestBody @Valid CredentialsDto credentials) {
        return ResponseEntity.ok(userService.login(credentials));
    }

    @PostMapping("/register")
    public ResponseEntity<@NonNull AuthResponse> register(@RequestBody @Valid SignUpDto userSignUp) {
        final AuthResponse response = userService.register(userSignUp);
        return ResponseEntity
                .created(URI.create("/users/" + response.user().id()))
                .body(response);
    }

    @PostMapping("/refresh")
    public ResponseEntity<@NonNull AuthResponse> refresh(@RequestBody @Valid RefreshTokenRequest request) {
        return ResponseEntity.ok(userService.refresh(request));
    }

    @PostMapping("/logout")
    public ResponseEntity<@NonNull Void> logout(@RequestBody @Valid RefreshTokenRequest request) {
        refreshTokenService.logout(request);

        return ResponseEntity.ok().build();
    }

    @PostMapping("/logout-all")
    public ResponseEntity<@NonNull Void> logoutAll(@AuthenticationPrincipal UserDto user) {
        if (user == null) {
            throw new UserException("Authentication required", HttpStatus.UNAUTHORIZED);
        }
        userService.logoutAll(user.email());

        return ResponseEntity.ok().build();
    }
}
