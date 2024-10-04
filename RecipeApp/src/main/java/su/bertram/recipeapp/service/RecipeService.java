package su.bertram.recipeapp.service;

import su.bertram.recipeapp.repository.RecipeRepository;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import su.bertram.recipeapp.kafka.producer.MessageProducer;

@Service
public class RecipeService {
    private final RecipeRepository recipeRepository;
    private final MessageProducer messageProducer;

    public RecipeService(RecipeRepository recipeRepository, MessageProducer messageProducer){
        this.recipeRepository = recipeRepository;
        this.messageProducer = messageProducer;
    }

    @Transactional
    public int deleteById(long id){
        // Delete from database.
        int deletedRecipeId = recipeRepository.deleteById(id);

        // Publish that the comment was deleted if it was successfully deleted.
        //if (deletedRecipeId != 0)
        messageProducer.sendMessage("delete-recipe-topic", String.format("{ \"id\": %s }", id));

        return deletedRecipeId;
    }
}
