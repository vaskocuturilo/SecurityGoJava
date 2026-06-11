package com.example.java.service;

import com.example.java.dto.Task;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class TaskService implements ITaskService {

    private final Map<String, String> tasks = new ConcurrentHashMap<>();


    @Override
    public Map<String, String> getTasks() {
        tasks.put("1", "Task1");
        tasks.put("2", "Task2");
        tasks.put("3", "Task3");

        return Collections.unmodifiableMap(tasks);
    }

    @Override
    public Task createTask(Task task) {
        if (tasks.containsKey(task.key())) {
            throw new IllegalArgumentException("Task with key already exists");
        }
        tasks.put(task.key(), task.value());

        return task;
    }
}
