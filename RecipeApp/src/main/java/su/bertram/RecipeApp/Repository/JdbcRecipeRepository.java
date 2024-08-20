package su.bertram.RecipeApp.Repository;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.IncorrectResultSizeDataAccessException;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import su.bertram.RecipeApp.Model.Recipe;

@Repository
public class JdbcRecipeRepository implements RecipeRepository{

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Override
    public int save(Recipe recipe){
        String cmd = "INSERT INTO recipe(title, url) VALUES(?, ?)";
        Object[] params = { recipe.getTitle(), recipe.getUrl()};
        return jdbcTemplate.update(cmd, params);
    }

    @Override
    public int update(Recipe recipe){
        String cmd = "UPDATE recipe SET title=?, url=? WHERE id=?";
        Object[] params = { recipe.getTitle(), recipe.getUrl(), recipe.getRecipeId()};
        return jdbcTemplate.update(cmd, params);
    }

    @Override
    public Recipe findById(long id){
        try {
            Recipe recipe = jdbcTemplate.queryForObject("SELECT recipe_id, title, url, created_at AS createdAt FROM recipe WHERE id=?",
                    BeanPropertyRowMapper.newInstance(Recipe.class), id);

            return recipe;
        } catch (IncorrectResultSizeDataAccessException e) {
            return null;
        }
    }

    @Override
    public int deleteById(long id) {
        return jdbcTemplate.update("DELETE FROM recipe WHERE id=?", id);
    }

    @Override
    public List<Recipe> findAll() {
        return jdbcTemplate.query("SELECT recipe_id, title, url, created_at AS createdAt FROM recipe", BeanPropertyRowMapper.newInstance(Recipe.class));
    }

    @Override
    public List<Recipe> findByTitleContaining(String title) {
        String q = "SELECT recipe_id, title, url, created_at AS createdAt FROM recipe WHERE title ILIKE '%" + title + "%'";

        return jdbcTemplate.query(q, BeanPropertyRowMapper.newInstance(Recipe.class));
    }
}
