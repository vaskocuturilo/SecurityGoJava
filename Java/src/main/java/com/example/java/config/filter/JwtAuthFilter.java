package com.example.java.config.filter;

import com.auth0.jwt.exceptions.JWTVerificationException;
import com.example.java.config.UserAuthenticationProvider;
import com.example.java.config.path.SecurityConstants;
import com.example.java.dto.UserDto;
import com.example.java.service.UserService;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.NonNull;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;
import org.springframework.web.servlet.HandlerExceptionResolver;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;

@Component
@RequiredArgsConstructor
public class JwtAuthFilter extends OncePerRequestFilter {

    private final UserAuthenticationProvider userAuthenticationProvider;
    private final UserService userService;
    private final HandlerExceptionResolver handlerExceptionResolver;

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        final String path = request.getRequestURI();
        return Arrays.asList(SecurityConstants.PUBLIC_ROUTES).contains(path);
    }

    @Override
    protected void doFilterInternal(@NonNull HttpServletRequest request,
                                    @NonNull HttpServletResponse response,
                                    @NonNull FilterChain filterChain)
            throws ServletException, IOException {

        final String header = request.getHeader(HttpHeaders.AUTHORIZATION);

        if (header == null) {
            handlerExceptionResolver.resolveException(request, response, null,
                    new JWTVerificationException("Missing authorization header"));
            return;
        }

        final String[] authElements = header.split(" ");

        if (authElements.length != 2 || !"Bearer".equals(authElements[0])) {
            SecurityContextHolder.clearContext();
            handlerExceptionResolver.resolveException(request, response, null,
                    new JWTVerificationException("Malformed authorization header"));
            return;
        }

        try {
            final String email = userAuthenticationProvider.extractEmail(authElements[1]);
            final UserDto user = userService.findByLogin(email);

            SecurityContextHolder.getContext().setAuthentication(
                    new UsernamePasswordAuthenticationToken(
                            user, null, Collections.emptyList()));

        } catch (JWTVerificationException exception) {
            SecurityContextHolder.clearContext();
            handlerExceptionResolver.resolveException(request, response, null, exception);
            return;
        }

        filterChain.doFilter(request, response);
    }
}