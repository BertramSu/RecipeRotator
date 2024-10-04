package su.bertram.recipeapp.repository;

import su.bertram.recipeapp.model.Recipe;
import java.util.List;

public interface RecipeRepository {
    int save(Recipe recipe);

    int update(Recipe recipe);

    Recipe findById(long id);

    int deleteById(long id);

    List<Recipe> findAll();

    List<Recipe> findByTitleContaining(String title);

    List<Recipe> findRecipesByTagId(long tagId);
}
