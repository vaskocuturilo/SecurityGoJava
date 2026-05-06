package com.example.java.rest;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/tasks")
public class RestTaskControllerV1 {

    private final Map<String, String> tasks = new HashMap<>();

    @GetMapping
    public ResponseEntity<Map<String, String>> getTasks() {
        tasks.put("1", "Task1");
        tasks.put("2", "Task2");
        tasks.put("3", "Task3");

        return ResponseEntity.ok().body(tasks);
    }
}
