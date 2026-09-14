document.addEventListener("DOMContentLoaded", () => {
    const postCreationForm = document.getElementById("create-post-form") || document.getElementById("createPostForm");
    const categoryCheckboxes = document.querySelectorAll(".category-checkbox");
    const maxAllowedCategories = 3;

    // 1. Enforce validation on Multi-Category selections natively
    // Captures the rule: "When registered users are creating a post they can associate one or more categories to it."
    if (categoryCheckboxes.length > 0) {
        categoryCheckboxes.forEach(checkbox => {
            checkbox.addEventListener("change", () => {
                const selectedCount = document.querySelectorAll(".category-checkbox:checked").length;

                // UX Guardrail: Prevent selecting too many categories if you want a clean scope
                if (selectedCount > maxAllowedCategories) {
                    checkbox.checked = false;
                    alert(`Architectural Constraint: You can only map a maximum of ${maxAllowedCategories} categories to a single post blueprint.`);
                }
            });
        });
    }

    // 2. Validate structural content fields before submitting them to the HTTP delivery layer
    if (postCreationForm) {
        postCreationForm.addEventListener("submit", (e) => {
            const titleElement = document.getElementById("post-title") || document.getElementById("postTitle");
            const contentElement = document.getElementById("post-content") || document.getElementById("postContent");
            const titleInput = titleElement.value.trim();
            const contentInput = contentElement.value.trim();
            const selectedCategories = document.querySelectorAll(".category-checkbox:checked").length;

            if (titleInput === "" || contentInput === "") {
                e.preventDefault();
                alert("Validation Failure: Post Title and Core Content parameters cannot be empty string segments.");
                return;
            }

            if (selectedCategories === 0) {
                e.preventDefault();
                alert("Validation Failure: You must map at least one subforum category tag onto your blueprint.");
                return;
            }

            console.log("Post payload successfully validated. Passing stream up to delivery routing...");
        });
    }
});
