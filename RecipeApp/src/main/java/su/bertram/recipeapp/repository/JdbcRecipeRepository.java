package su.bertram.recipeapp.repository;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.IncorrectResultSizeDataAccessException;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import su.bertram.recipeapp.model.Recipe;

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
        String cmd = "UPDATE recipe SET title=?, url=? WHERE recipe_id=?";
        Object[] params = { recipe.getTitle(), recipe.getUrl(), recipe.getRecipeId()};
        return jdbcTemplate.update(cmd, params);
    }

    @Override
    public Recipe findById(long id){
        try {
            Recipe recipe = jdbcTemplate.queryForObject("SELECT * FROM recipe WHERE recipe_id=?",
                    BeanPropertyRowMapper.newInstance(Recipe.class), id);

            return recipe;
        } catch (IncorrectResultSizeDataAccessException e) {
            return null;
        }
    }

    @Override
    public int deleteById(long id) {
        return jdbcTemplate.update("DELETE FROM recipe WHERE recipe_id=?", id);
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

    @Override
    public List<Recipe> findRecipesByTagId(long tagId){
        String q = "SELECT\n" +
                "    r.recipe_id ,\n" +
                "    r.title,\n" +
                "    r.url\n" +
                "FROM \n" +
                "    recipetag AS rt \n" +
                "    INNER JOIN recipe AS r on rt.recipe_id = r.recipe_id\n" +
                "WHERE \n" +
                "    rt.tag_id = " +tagId;

        return jdbcTemplate.query(q, BeanPropertyRowMapper.newInstance(Recipe.class));
    }

}
