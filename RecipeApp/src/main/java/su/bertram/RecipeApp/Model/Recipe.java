package su.bertram.RecipeApp.Model;

import java.time.LocalDateTime;

public class Recipe {
    private long id;
    private String title;
    private String url;
    private LocalDateTime createdAt;

    // Constructors
    public Recipe() {
        // Default constructor
    }

    public Recipe(String title, String url) {
        this.title = title;
        this.url = url;
    }

    public Recipe(long recipeId, String title, String url) {
        this.id = recipeId;
        this.title = title;
        this.url = url;
    }

    // Getters and Setters
    public long getRecipeId() {
        return id;
    }

    public void setRecipeId(long recipeId) {
        this.id = recipeId;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url;
    }

    public LocalDateTime getCreatedAt() {
        return createdAt;
    }

    // toString method
    @Override
    public String toString() {
        return "Recipe{" +
                "recipeId=" + id +
                ", title='" + title + '\'' +
                ", url='" + url + '\'' +
                ", createdAt=" + createdAt +
                '}';
    }
}
