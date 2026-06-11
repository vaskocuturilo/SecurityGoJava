package com.example.java.service;

import com.example.java.dto.Task;

import java.util.Map;

public interface ITaskService {

    Map<String, String> getTasks();

    Task createTask(Task task);
}
