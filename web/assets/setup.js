// Copyright (C) 2023 David Sugar, Tycho Softworks
// This code is licensed under MIT license

function validateForm() {
    var input = document.getElementById("admin");
    input.value = input.value.trim();
    if (input.value === "") {
        alert("Please enter an admin username.");
        return false;
    }

    input = document.GetElementById("pass");
    var passwd = input.value
    if (passwd === "") {
        alert("Please enter a password.");
        return false;
    }

    input = document.GetElementById("verify");
    var verify = input.value
    if (verify !== passwd) {
        alert("Password does not match verify.");
        return false;
    }
}

