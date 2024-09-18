package su.bertram.RecipeApp.Controller;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import su.bertram.RecipeApp.Model.Recipe;
import su.bertram.RecipeApp.Repository.RecipeRepository;
import su.bertram.RecipeApp.Service.RecipeService;

@CrossOrigin(origins = "http://localhost:8081")
@RestController
@RequestMapping("/api")
public class RecipeController {

    @Autowired
    RecipeRepository recipeRepository;

    @Autowired
    RecipeService recipeService;

    @GetMapping("/recipes")
    public ResponseEntity<List<Recipe>> getAllRecipes(@RequestParam(required = false) String title){
        try{
            List<Recipe> recipes = new ArrayList<Recipe>();

            if (title == null)
                recipes.addAll(recipeRepository.findAll());
            else
                recipes.addAll(recipeRepository.findByTitleContaining(title));

            if (recipes.isEmpty())
                return new ResponseEntity<>(HttpStatus.NO_CONTENT);

            return new ResponseEntity<>(recipes, HttpStatus.OK);
        } catch (Exception e){
            return new ResponseEntity<>(null, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @GetMapping("/recipe/{id}")
    public ResponseEntity<Recipe> getRecipeById(@PathVariable("id") long id){
        Recipe recipe = recipeRepository.findById(id);

        if (recipe != null)
            return new ResponseEntity<>(recipe, HttpStatus.OK);
        else
            return new ResponseEntity<>(HttpStatus.NO_CONTENT);
    }

    @PostMapping("/recipe")
    public ResponseEntity<String> createRecipe(@RequestBody Recipe recipe){
        try {
            recipeRepository.save(new Recipe(recipe.getTitle(), recipe.getUrl()));
            return new ResponseEntity<>("Recipe was created successfully.", HttpStatus.CREATED);
        } catch (Exception e) {
            return new ResponseEntity<>(null, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @PutMapping("/recipe/{id}")
    public ResponseEntity<String> updateRecipe(@PathVariable("id") long id, @RequestBody Recipe recipe){
        Recipe _recipe = recipeRepository.findById(id);

        if (_recipe != null){
            _recipe.setTitle(recipe.getTitle());
            _recipe.setUrl(recipe.getUrl());

            recipeRepository.update(_recipe);
            return new ResponseEntity<>("Recipe was successfully updated.", HttpStatus.OK);
        }else
            return new ResponseEntity<>("Cannot find Recipe with id=" + id, HttpStatus.OK);
    }

    @DeleteMapping("recipe/{id}")
    public ResponseEntity<String> deleteRecipe(@PathVariable("id") long id){
        try {
            int deletedId = recipeService.deleteById(id);
            if (deletedId == 0)
                return new ResponseEntity<>("Cannot find Recipe with id=" + id, HttpStatus.OK);

            return new ResponseEntity<>("Recipe was deleted successfully.", HttpStatus.OK);
        } catch (Exception e){
            return new ResponseEntity<>("Cannot delete recipe.", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}
