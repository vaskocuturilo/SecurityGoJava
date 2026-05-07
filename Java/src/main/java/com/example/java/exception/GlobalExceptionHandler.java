package com.example.java.exception;

import com.auth0.jwt.exceptions.JWTVerificationException;
import com.auth0.jwt.exceptions.TokenExpiredException;
import jakarta.persistence.EntityNotFoundException;
import lombok.NonNull;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.time.LocalDateTime;
import java.util.stream.Collectors;

@RestControllerAdvice
public class GlobalExceptionHandler {

    @ExceptionHandler(EntityNotFoundException.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleEntityNotFound(EntityNotFoundException exception) {
        final ErrorResponseDto issue = new ErrorResponseDto(
                "The user issue",
                HttpStatus.NOT_FOUND.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.NOT_FOUND)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }

    @ExceptionHandler(UserException.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleUserExist(UserException exception) {
        final ErrorResponseDto issue = new ErrorResponseDto(
                "The user issue",
                HttpStatus.BAD_REQUEST.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }

    @ExceptionHandler(TokenExpiredException.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleTokenExpired(TokenExpiredException exception) {
        final ErrorResponseDto issue = new ErrorResponseDto(
                "Token has expired",
                HttpStatus.UNAUTHORIZED.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.UNAUTHORIZED)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }

    @ExceptionHandler(JWTVerificationException.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleInvalidToken(JWTVerificationException exception) {
        final ErrorResponseDto issue = new ErrorResponseDto(
                "Invalid token",
                HttpStatus.UNAUTHORIZED.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.UNAUTHORIZED)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }

    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleIllegalArgument(IllegalArgumentException exception) {
        final ErrorResponseDto issue = new ErrorResponseDto(
                "Invalid argument",
                HttpStatus.BAD_REQUEST.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleNotValidArgument(MethodArgumentNotValidException exception) {
        final String message = exception.getBindingResult()
                .getFieldErrors()
                .stream()
                .map(FieldError::getDefaultMessage)
                .collect(Collectors.joining(", "));

        final ErrorResponseDto issue = new ErrorResponseDto(
                "Invalid argument",
                HttpStatus.BAD_REQUEST.value(),
                message,
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<@NonNull ErrorResponseDto> handleGeneric(Exception exception) {
        final ErrorResponseDto issue = new ErrorResponseDto(
                "Unexpected error:",
                HttpStatus.INTERNAL_SERVER_ERROR.value(),
                exception.getMessage(),
                LocalDateTime.now());

        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .contentType(MediaType.APPLICATION_JSON)
                .body(issue);
    }
}
