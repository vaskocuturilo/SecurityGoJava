package com.example.java.rest;

import com.example.java.annotation.RateLimiter;
import com.example.java.dto.Task;
import com.example.java.service.ITaskService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api/v1/tasks")
@RequiredArgsConstructor
public class RestTaskControllerV1 {

    private final ITaskService taskService;

    @GetMapping
    @RateLimiter(key = "get tasks")
    public ResponseEntity<Map<String, String>> getTasks() {
        return ResponseEntity.ok().body(taskService.getTasks());
    }

    @PostMapping
    @RateLimiter(key = "create task")
    public ResponseEntity<Task> createTask(@Valid @RequestBody Task task) {
        final Task createdTask = taskService.createTask(task);
        return ResponseEntity.status(HttpStatus.CREATED).body(createdTask);
    }
}
