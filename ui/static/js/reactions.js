document.addEventListener("DOMContentLoaded", () => {
    const reactionButtons = document.querySelectorAll(".reaction-btn[data-target-id]");

    reactionButtons.forEach((button) => {
        button.addEventListener("click", async () => {
            const targetButtons = [...reactionButtons].filter((candidate) => (
                candidate.dataset.targetId === button.dataset.targetId &&
                candidate.dataset.targetType === button.dataset.targetType
            ));
            const previousState = targetButtons.map((candidate) => ({
                button: candidate,
                active: candidate.dataset.active === "true",
                count: getCount(candidate),
            }));

            const wasActive = button.dataset.active === "true";
            targetButtons.forEach((candidate) => {
                const previous = previousState.find((state) => state.button === candidate);
                const isSelected = candidate === button && !wasActive;
                candidate.dataset.active = String(isSelected);
                if (isSelected) {
                    setCount(candidate, getCount(candidate) + 1);
                } else if (previous.active && (!isSelected || candidate === button)) {
                    setCount(candidate, Math.max(0, getCount(candidate) - 1));
                }
            });

            targetButtons.forEach((candidate) => { candidate.disabled = true; });
            try {
                const response = await fetch("/interaction", {
                    method: "POST",
                    credentials: "same-origin",
                    headers: { "Content-Type": "application/x-www-form-urlencoded" },
                    body: new URLSearchParams({
                        target_id: button.dataset.targetId,
                        target_type: button.dataset.targetType,
                        value: button.dataset.value,
                    }),
                });
                if (response.status === 401) {
                    throw new Error("You must be logged in to like or dislike.");
                }
                if (!response.ok) {
                    throw new Error("Unable to save your reaction.");
                }
            } catch (error) {
                previousState.forEach((state) => {
                    state.button.dataset.active = String(state.active);
                    setCount(state.button, state.count);
                });
                showError(button.closest(".post-reaction-bar, .comment-reactions"), error.message);
            } finally {
                targetButtons.forEach((candidate) => { candidate.disabled = false; });
            }
        });
    });

    function getCount(button) {
        return Number(button.querySelector(".reaction-count")?.textContent || 0);
    }

    function setCount(button, count) {
        const countElement = button.querySelector(".reaction-count");
        if (countElement) {
            countElement.textContent = String(count);
        }
    }

    function showError(container, message) {
        const errorElement = container?.querySelector(".interaction-error");
        if (!errorElement) return;
        errorElement.textContent = message;
        window.setTimeout(() => {
            errorElement.textContent = "";
        }, 4000);
    }
});
