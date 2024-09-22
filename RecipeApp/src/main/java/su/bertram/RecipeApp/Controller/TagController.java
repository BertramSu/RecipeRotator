package su.bertram.RecipeApp.Controller;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;

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
import su.bertram.RecipeApp.Model.Tag;
import su.bertram.RecipeApp.Repository.TagRepository;

@CrossOrigin(origins = "http://localhost:8081")
@RestController
@RequestMapping("/api")
public class TagController {
    @Autowired
    TagRepository tagRepository;

    @GetMapping("/tags")
    public ResponseEntity<List<Tag>> getAllRecipes(@RequestParam(required = false) String title){
        try{
            List<Tag> tags = new ArrayList<Tag>();
            tags.addAll(tagRepository.findAll());

            if (tags.isEmpty())
                return new ResponseEntity<>(HttpStatus.NO_CONTENT);

            return new ResponseEntity<>(tags, HttpStatus.OK);
        } catch (Exception e){
            return new ResponseEntity<>(null, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @GetMapping("/tag/{id}")
    public ResponseEntity<Optional<Tag>> getRecipeById(@PathVariable("id") long id){
        Optional<Tag> tag = tagRepository.findById(id);

        if (tag.isPresent())
            return new ResponseEntity<>(tag, HttpStatus.OK);
        else
            return new ResponseEntity<>(HttpStatus.NO_CONTENT);
    }

    @PostMapping("/tag")
    public ResponseEntity<String> createRecipe(@RequestBody Tag tag){
        try {
            tagRepository.save(new Tag(tag.getName()));
            return new ResponseEntity<>("Recipe was created successfully.", HttpStatus.CREATED);
        } catch (Exception e) {
            return new ResponseEntity<>(null, HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }

    @PutMapping("/tag/{id}")
    public ResponseEntity<String> updateRecipe(@PathVariable("id") long id, @RequestBody Tag tag){
        Optional<Tag> _tag = tagRepository.findById(id);

        if (_tag.isPresent()){
            _tag.get().setName(tag.getName());
            tagRepository.save(_tag.get());
            return new ResponseEntity<>("Recipe was successfully updated.", HttpStatus.OK);
        }else
            return new ResponseEntity<>("Cannot find Recipe with id=" + id, HttpStatus.OK);
    }

    @DeleteMapping("tag/{id}")
    public ResponseEntity<String> deleteRecipe(@PathVariable("id") long id){
        try {
            Optional<Tag> tagToBeDeleted = tagRepository.findById(id);
            if (tagToBeDeleted.isEmpty())
                return new ResponseEntity<>("Cannot find Recipe with id=" + id, HttpStatus.OK);

            tagRepository.deleteById(id);
            return new ResponseEntity<>("Recipe was deleted successfully.", HttpStatus.OK);
        } catch (Exception e){
            return new ResponseEntity<>("Cannot delete recipe.", HttpStatus.INTERNAL_SERVER_ERROR);
        }
    }
}
