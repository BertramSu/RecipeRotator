package su.bertram.RecipeApp.Repository;

import su.bertram.RecipeApp.Model.Recipe;
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
